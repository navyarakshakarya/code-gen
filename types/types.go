package types

// ProjectInfo contains all analyzed project information
type ProjectInfo struct {
	ModuleName  string
	PackageName string
	ProjectDir  string
	Interfaces  map[string]*InterfaceInfo
	Structs     map[string]*StructInfo
	Imports     map[string]string // package -> import path
}

// InterfaceInfo represents an analyzed interface
type InterfaceInfo struct {
	Name              string
	Package           string
	FilePath          string
	Methods           []MethodInfo
	Layer             LayerType
	RelatedInterfaces []string
	Comments          []string
}

// StructInfo represents an analyzed struct
type StructInfo struct {
	Name     string
	Package  string
	FilePath string
	Fields   []FieldInfo
	Comments []string
}

// MethodInfo represents a method in an interface
type MethodInfo struct {
	Name       string
	Params     []ParamInfo
	Returns    []ParamInfo
	HasContext bool
	HasError   bool
	Comments   []string
}

// ParamInfo represents a parameter or return value
type ParamInfo struct {
	Name string
	Type string
}

// FieldInfo represents a struct field
type FieldInfo struct {
	Name     string
	Type     string
	Tag      string
	Embedded bool
}

// LayerType represents the architectural layer
type LayerType string

const (
	RepositoryLayer LayerType = "repository"
	UseCaseLayer    LayerType = "usecase"
	HandlerLayer    LayerType = "handler"
	ServiceLayer    LayerType = "service"
)

// String returns the string representation of LayerType
func (l LayerType) String() string {
	return string(l)
}
