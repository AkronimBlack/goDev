package templates

var depMap = map[string]func() []byte{
	"gin": GinDep,
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
