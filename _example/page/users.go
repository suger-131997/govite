package page

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/suger-131997/govite"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	JoinedAt string `json:"joinedAt"`
}

type UsersProps struct {
	Users       []User `json:"users"`
	CurrentPage int    `json:"currentPage"`
	TotalPages  int    `json:"totalPages"`
	TotalUsers  int    `json:"totalUsers"`
	PageSize    int    `json:"pageSize"`
}

var dummyUsers = func() []User {
	names := []string{
		"Alice Johnson", "Bob Smith", "Carol Williams", "David Brown", "Eve Davis",
		"Frank Miller", "Grace Wilson", "Henry Moore", "Ivy Taylor", "Jack Anderson",
		"Karen Thomas", "Liam Jackson", "Mia White", "Noah Harris", "Olivia Martin",
		"Peter Garcia", "Quinn Martinez", "Rose Robinson", "Sam Clark", "Tina Lewis",
		"Uma Lee", "Victor Walker", "Wendy Hall", "Xavier Allen", "Yara Young",
		"Zane Hernandez", "Amy King", "Brian Wright", "Cathy Lopez", "Daniel Hill",
		"Eleanor Scott", "Felix Green", "Gina Adams", "Hank Baker", "Iris Nelson",
		"Jake Carter", "Kelly Mitchell", "Leo Perez", "Megan Roberts", "Neil Turner",
		"Oscar Phillips", "Paula Campbell", "Quinn Parker", "Rachel Evans", "Steve Edwards",
		"Tara Collins", "Ulrich Stewart", "Violet Sanchez", "Wayne Morris", "Xena Rogers",
	}
	roles := []string{"Admin", "Editor", "Viewer", "Developer", "Manager"}
	statuses := []string{"Active", "Inactive", "Pending"}

	users := make([]User, len(names))
	base := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i, name := range names {
		users[i] = User{
			ID:       i + 1,
			Name:     name,
			Email:    generateEmail(name),
			Role:     roles[i%len(roles)],
			Status:   statuses[i%len(statuses)],
			JoinedAt: base.AddDate(0, i, i*3).Format("2006-01-02"),
		}
	}
	return users
}()

func generateEmail(name string) string {
	email := ""
	lower := false
	for _, ch := range name {
		if ch == ' ' {
			email += "."
			lower = true
		} else {
			if lower || len(email) == 0 {
				if ch >= 'A' && ch <= 'Z' {
					email += string(ch + 32)
				} else {
					email += string(ch)
				}
				lower = false
			} else {
				if ch >= 'A' && ch <= 'Z' {
					email += string(ch + 32)
				} else {
					email += string(ch)
				}
			}
		}
	}
	return email + "@example.com"
}

func NewUsersHandler() *govite.PageHandler[UsersProps] {
	return govite.NewPageHandler[UsersProps](govite.PageHandlerConfig[UsersProps]{
		EntryPoint: "page/users.tsx",
		HandleFunc: func(r *http.Request, render func(ctx context.Context, props UsersProps)) {
			ctx := r.Context()
			ctx = govite.WithTitle(ctx, "Go + Vite Demo: Users")

			const defaultPageSize = 10

			page := 1
			if p := r.URL.Query().Get("page"); p != "" {
				if n, err := strconv.Atoi(p); err == nil && n > 0 {
					page = n
				}
			}

			pageSize := defaultPageSize
			if s := r.URL.Query().Get("pageSize"); s != "" {
				if n, err := strconv.Atoi(s); err == nil && n > 0 && n <= 50 {
					pageSize = n
				}
			}

			total := len(dummyUsers)
			totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
			if page > totalPages {
				page = totalPages
			}

			start := (page - 1) * pageSize
			end := start + pageSize
			if end > total {
				end = total
			}

			render(ctx, UsersProps{
				Users:       dummyUsers[start:end],
				CurrentPage: page,
				TotalPages:  totalPages,
				TotalUsers:  total,
				PageSize:    pageSize,
			})
		},
	})
}
