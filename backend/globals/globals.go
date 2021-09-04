package globals

import (
	"github.com/maxkruse/magnusopus/backend/structs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Logger *logrus.Logger
	Config structs.Config
	DBConn *gorm.DB
)
