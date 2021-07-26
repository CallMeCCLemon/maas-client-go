package maasclient

import (
	"context"
	"encoding/json"
)

type Spaces interface {
	List(ctx context.Context) ([]Space, error)
}

type Space interface {
	Name() string
	Subnets() []Subnet
}

type space struct {
	name    string
	subnets []*subnet
}

func (s *space) Name() string {
	return s.name
}

func (s *space) Subnets() []Subnet {
	return subnetSliceToInterface(s.subnets)
}

func (s *space) UnmarshalJSON(data []byte) error {
	des := &struct {
		Name    string    `json:"name"`
		Subnets []*subnet `json:"subnets"`
	}{}

	err := json.Unmarshal(data, des)
	if err != nil {
		return err
	}

	s.name = des.Name
	s.subnets = des.Subnets

	return nil
}

type spaces struct {
	Controller
}

func (ss *spaces) List(ctx context.Context) ([]Space, error) {
	res, err := ss.client.Get(ctx, ss.apiPath, ss.params.Values())
	if err != nil {
		return nil, err
	}

	var obj []*space
	err = unMarshalJson(res, &obj)
	if err != nil {
		return nil, err
	}

	return spaceStructSliceToInterface(obj, ss.client), nil
}

func spaceStructSliceToInterface(in []*space, client Client) []Space {
	var out []Space
	for _, s := range in {
		out = append(out, spaceStructToInterface(s, client))
	}
	return out
}

func spaceStructToInterface(in *space, client Client) Space {
	return in
}

func NewSpacesClient(client *authenticatedClient) Spaces {
	return &spaces{
		Controller: Controller{
			client:  client,
			apiPath: "/spaces/",
			params:  ParamsBuilder(),
		},
	}
}
