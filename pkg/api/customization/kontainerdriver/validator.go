package kontainerdriver

import (
	"fmt"
	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/labels"
)

type Validator struct {
	KontainerDriverLister v3.KontainerDriverLister
}

func (v *Validator) Validator(request *types.APIContext, schema *types.Schema, data map[string]interface{}) error {
	var spec v3.KontainerDriverSpec

	if err := convert.ToObj(data, &spec); err != nil {
		return httperror.WrapAPIError(err, httperror.InvalidBodyContent, "Kontainer driver spec conversion error")
	}

	return v.validateKontainerDriverURL(request, spec)
}

func (v *Validator) validateKontainerDriverURL(request *types.APIContext, spec v3.KontainerDriverSpec) error {
	kontainerDrivers, err := v.KontainerDriverLister.List("", labels.NewSelector())
	if err != nil {
		return httperror.WrapAPIError(err, httperror.ServerError, "Failed to list kontainer drivers")
	}

	for _, driver := range kontainerDrivers {
		if driver.Spec.URL == spec.URL && driver.Name != request.ID {
			return httperror.NewAPIError(httperror.Conflict, fmt.Sprintf("Driver URL already in use: %s", spec.URL))
		}
	}

	return nil
}
