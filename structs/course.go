package structs

type Course struct {
	DeptKey     string
	Key         string
	Name        string
	Desc        string
	ProfileLink string
	Credits     int
	Aliases     []string
	ReqProps    RequisiteProps
}

/*
* requisite properteis include:
* 		enforced
* 		advisory
* 		standing
* 		permission of instructor
 */
type RequisiteProps struct {
	None            bool
	Enforced        []RequisiteCourse
	Advisory        []RequisiteCourse
	Standing        []Standing
	InstructorPerms bool
	Notes           string
	Raw             string
}

// can be advisory or required
type Standing struct {
	Name string
	Kind string
}

type RequisiteCourse struct {
	Key               string
	OrEquivalent      bool
	KnownEquivalences []string
	canAccompany      bool
}
