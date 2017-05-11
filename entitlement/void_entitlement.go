package entitlement

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/docker/libentitlement/context"
	"github.com/docker/libentitlement/parser"
	"strings"
)

type VoidEntitlementEnforceCallback func(*context.Context) (*context.Context, error)

type VoidEntitlement struct {
	domain           []string
	id               string
	enforce_callback VoidEntitlementEnforceCallback
}

func NewVoidEntitlement(fullName string, callback VoidEntitlementEnforceCallback) *VoidEntitlement {
	domain, id, err := parser.ParseVoidEntitlement(fullName)
	if err != nil {
		logrus.Errorf("Couldn't not create entitlement for %v\n", fullName)
		return nil
	}

	return &VoidEntitlement{domain: domain, id: id, enforce_callback: callback}
}

func (e *VoidEntitlement) Domain() (string, error) {
	if len(e.domain) < 1 {
		id, err := e.Identifier()
		if err != nil {
			return "", fmt.Errorf("No domain or id found for current entitlement")
		}

		return "", fmt.Errorf("No domain found for entitlement %s", id)
	}

	return strings.Join(e.domain, "."), nil
}

func (e *VoidEntitlement) Identifier() (string, error) {
	if e.id == "" {
		return "", fmt.Errorf("No identifier found for current entitlement")
	}

	return e.id, nil
}

// Value() should not be called on a void entitlement
func (e *VoidEntitlement) Value() (string, error) {
	return "", nil
}

func (e *VoidEntitlement) Enforce(ctx *context.Context) (*context.Context, error) {
	domain, _ := e.Domain()
	id, _ := e.Identifier()

	if e.enforce_callback == nil {
		return nil, fmt.Errorf("Invalid enforcement callback for entitlement %v.%v", domain, id)
	}

	newContext, err := e.enforce_callback(ctx)
	if err != nil {
		return nil, err
	}

	return newContext, err
}