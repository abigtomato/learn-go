package designpattern

import (
	"fmt"
	"testing"
)

// 命令的接收者-医生，命令具体的执行者
type Doctor struct{}

func (d Doctor) treatEye() {
	fmt.Println("医生治疗眼睛")
}

func (d Doctor) treatNose() {
	fmt.Println("医生治疗鼻子")
}

// 抽象的命令
type Command interface {
	// 治疗行为
	Treat()
}

// 具体的命令-治疗眼睛的病单
type CommandTreatEye struct {
	doctor *Doctor
}

func (cte CommandTreatEye) Treat() {
	cte.doctor.treatEye()
}

// 具体的命令-治疗鼻子的病单
type CommandTreatNose struct {
	doctor *Doctor
}

func (ctn CommandTreatNose) Treat() {
	ctn.doctor.treatNose()
}

// 命令的使用者-护士，调动命令的一方
type Nurse struct {
	// 命令列表
	CmdList []Command
}

// 调用所有命令
func (n Nurse) Notify() {
	if n.CmdList == nil {
		return
	}
	for _, cmd := range n.CmdList {
		cmd.Treat()
	}
}

func TestCommand(t *testing.T) {
	doctor := new(Doctor)
	cmdEye := CommandTreatEye{doctor}
	cmdNose := CommandTreatNose{doctor}

	nurse := new(Nurse)
	nurse.CmdList = append(nurse.CmdList, &cmdEye)
	nurse.CmdList = append(nurse.CmdList, &cmdNose)
	nurse.Notify()
}
