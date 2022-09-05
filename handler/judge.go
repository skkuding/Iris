package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cranemont/judge-manager/file"
	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/ingress/rmq"
	"github.com/cranemont/judge-manager/sandbox"
)

var handler = "JudgeHandler"

type JudgeResult struct {
	StatusCode int                   `json:"statusCode"` // handler's status code
	Data       judge.JudgeTaskResult `json:"data"`
}

type JudgeHandler struct {
	judger *judge.Judger
	config *sandbox.LanguageConfig
}

func NewJudgeHandler(
	judger *judge.Judger,
	config *sandbox.LanguageConfig,
) *JudgeHandler {
	return &JudgeHandler{
		judger: judger,
		config: config,
	}
}

// handle top layer logical flow
func (h *JudgeHandler) Handle(request rmq.JudgeRequest) (result JudgeResult, err error) {
	res := JudgeResult{StatusCode: INTERNAL_SERVER_ERROR, Data: judge.JudgeTaskResult{}}
	task := judge.NewTask(request)
	task.StartedAt = time.Now()
	dir := task.GetDir()

	defer func() {
		file.RemoveDir(task.GetDir())
		fmt.Println(time.Since(task.StartedAt)) // for debug
	}()

	if err := file.CreateDir(dir); err != nil {
		return res, fmt.Errorf("%s: failed to create directory: %w", handler, err)
	}

	srcPath, err := h.config.MakeSrcPath(dir, task.GetLanguage())
	if err != nil {
		return res, fmt.Errorf("%s: failed to create src path: %w", handler, err)
	}
	if err := file.CreateFile(srcPath, task.GetCode()); err != nil {
		return res, fmt.Errorf("%s: failed to create src file: %w", handler, err)
	}

	err = h.judger.Judge(task)
	if err != nil {
		if errors.Is(err, judge.ErrTestcaseGet) {
			res.StatusCode = TESTCASE_GET_FAILED
		} else if !errors.Is(err, judge.ErrCompile) {
			return res, fmt.Errorf("%s: judge failed: %w", handler, err)
		}
		res.StatusCode = COMPILE_ERROR
	} else {
		res.StatusCode = SUCCESS
	}

	res.Data = task.Result
	fmt.Println("JudgeHandler: Handle Done!")
	return res, nil
}

func (h *JudgeHandler) ResultToJson(result JudgeResult) string {
	res, err := json.Marshal(result)
	if err != nil {
		// 적절한 err 처리
		panic(err)
	}
	return string(res)
}
