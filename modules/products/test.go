package products

// package products

// // Component
// type Service struct {
// 	// it will likely have some dependencies on other components
// 	CommentRepo CommentRepo
// }

// // these dependencies should be modelled as interfaces to help ensure
// // our components are loosely coupled.
// type CommentRepo interface {
// 	GetComments() ([]Comment, error)
// }

// // Comment - in this example, imagine our component is doing
// // stuff processing of comments on this website.
// type Comment struct {
// 	Author string
// 	Body   string
// 	Slug   string
// 	*Service
// }

// // NewService - our constructor function
// func NewService(cmtRepo CommentRepo) (*Service, error) {
// 	svc := &Service{
// 		CommentRepo: cmtRepo,
// 	}
// 	// handles other potentially more complex setup logic
// 	// for our component, there could be calls to downstream
// 	// dependencies to check connections etc that could return
// 	// errors
// 	return svc, nil
// }

// // DoesStuff - a method that takes a pointer receiver to an
// // instantiated Component
// func (c *Comment) DoesStuff() error {
// 	comments, err := c.
// 	if err != nil {
// 		return err
// 	}
// 	// do additional things with the returned comments
// }
