package gotrus

var CtxVar string = ""

type (

	// Obj is an interface describing an object, useful when retrieving properties
	Obj interface {
		// Guid returns the guid (uuidv4) of this object.
		Guid() string

		// StringProp retrieves a named string property. If not yet loaded, it gets retrieved and returned.
		StringProp(n string) (string, error)

		// IntProp retrieves a named integer (int64) property. If not yet loaded, it gets retrieved and returned.
		IntProp(n string) (int64, error)

		// BoolProp retrieves a named boolean property. If not yet loaded, it gets retrieved and returned.
		BoolProp(n string) (bool, error)

		// FloatProp retrieves a named float (float64) property. If not yet loaded, it gets retrieved and returned.
		FloatProp(n string) (float64, error)
	}

	// Petrus is the interface describing an entire Petrus object.
	Petrus interface {
		// User returns the decoded user.
		User() Obj

		// Company returns the decoded user's company.
		Company(string) Obj

		// Companies returns a slice of all companies.
		Companies() []Obj

		// CompaniesByGuid returns a map of all companies, keyed by company guid.
		CompaniesByGuid() map[string]Obj

		// Whitelabel returns the decoded user's whitelabel.
		Whitelabel() Obj

		// Token returns the token of this petrus instance, the JTI claim.
		Token() string

		// Claim returns a single claim of the JWT
		Claim(claim string) (interface{}, bool)

		// HasRole returns true iff the Petrus user has a specific role.
		HasRole(company string, role string) bool

		// HasOneRole returns true iff the Petrus user has one of a set of roles.
		HasOneRole(company string, roles ...string) bool

		// HasAllRoles returns true iff the Petrus user has all roles of a set of roles.
		HasAllRoles(company string, roles ...string) bool
	}
)
