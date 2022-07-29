package setup

import "github.com/lamoda/gonkey/models"

type StartUpFunction func(testInterface models.TestInterface) error

type TeardownFunction func(testInterface models.TestInterface) error
