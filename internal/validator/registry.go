package validator

// Operation type IDs, mirrored from the operation_types seed data. Kept here so
// the rule wiring reads clearly; the source of truth remains the database.
const (
	opNormalPurchase      = 1
	opInstallmentPurchase = 2
	opWithdrawal          = 3
	opCreditVoucher       = 4
)

// Registry maps an operation type to the ordered set of rules it must satisfy.
// Adding a new operation type is a single new entry; adding a rule to a group of
// types touches only the affected entries.
type Registry struct {
	rules map[int][]TransactionValidator
}

// NewRegistry wires the default rule set. Every operation type requires a
// positive amount; rules that differ per type are added to the relevant entries.
func NewRegistry() *Registry {
	positive := PositiveAmount{}

	return &Registry{
		rules: map[int][]TransactionValidator{
			opNormalPurchase:      {positive},
			opInstallmentPurchase: {positive},
			opWithdrawal:          {positive},
			opCreditVoucher:       {positive},
		},
	}
}

// Validate runs every rule configured for the operation type, returning the
// first failure. An operation type with no configured rules is treated as having
// no extra constraints; existence of the type is checked earlier by the service.
func (r *Registry) Validate(operationTypeID int, in Input) error {
	for _, rule := range r.rules[operationTypeID] {
		if err := rule.Validate(in); err != nil {
			return err
		}
	}
	return nil
}

// compile-time assurance that the concrete rules satisfy the interface.
var _ TransactionValidator = PositiveAmount{}
