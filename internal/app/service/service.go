package service

type Service struct {
}

type Post struct {
	Author string `json:"author"`
	Subs   int    `json:"subs"`
	Date   string `json:"date"`
	Text   string `json:"text"`
	ImgUrl string `json:"imgUrl"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

func New() *Service {
	return &Service{}
}

func (s *Service) GenerateNewPosts() Posts {
	posts := Posts{
		Posts: []Post{
			{
				Author: "John Doe",
				Subs:   100,
				Date:   "2024-12-27",
				Text:   "This is a sample post",
				ImgUrl: "http://example.com/image.jpg",
			},
			{
				Author: "Jane Doe",
				Subs:   200,
				Date:   "2024-12-26",
				Text:   "Another example post",
				ImgUrl: "http://example.com/image2.jpg",
			},
		},
	}
	return posts
}
