package templates

var depMap = map[string]func() []byte{
	"gin":    GinDep,
	"mysql":  GormDep,
	"logrus": LogrusDep,
}

func GetDependency(item string) []byte {
	if x, ok := depMap[item]; ok {
		return x()
	}
	return nil
}

/*GinDep exert from the go.mod file for gin*/
func GinDep() []byte {
	return []byte(`    github.com/gin-contrib/cors v1.3.1
    github.com/gin-gonic/gin v1.6.3`)
}

func GormDep() []byte {
	return []byte(`    gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.6`)
}

func LogrusDep() []byte {
	return []byte(`    github.com/sirupsen/logrus v1.8.1`)
}
