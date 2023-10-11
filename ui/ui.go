package ui

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/td0m/poc-doorman/entities"
	"github.com/td0m/poc-doorman/u"
	"github.com/td0m/poc-doorman/ui/templates"
)

type Template struct {
	Primary   string
	DependsOn []string
}

func t(p string, rest ...string) Template {
	return Template{
		Primary:   p,
		DependsOn: append([]string{"layout.html"}, rest...),
	}
}

func (t Template) Build() *template.Template {
	all := append([]string{t.Primary}, t.DependsOn...)
	funcs := map[string]any{
		"add": func(a, b int) int {
			return a + b
		},
	}
	return template.Must(template.New("").Funcs(funcs).ParseFS(templates.FS, all...))
}

type Ctx struct {
	context.Context
	*http.Request
	http.ResponseWriter

	Data   map[string]any
	Errors struct {
		Internal error
	}
}

// func (c *Ctx) Context()

func (t Template) Serve(f func(ctx *Ctx) error) http.HandlerFunc {
	tmpl := t.Build()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Ctx{
			Context:        r.Context(),
			Data:           map[string]any{},
			Request:        r,
			ResponseWriter: w,
		}

		if err := f(ctx); err != nil {
			ctx.Errors.Internal = err
		}

		name := t.Primary
		if hxRequest := r.Header.Get("HX-Request"); hxRequest == "true" {
			name = "hx"
		}
		if hxTarget := r.Header.Get("HX-Target"); hxTarget != "" {
			name = hxTarget
		}

		err := tmpl.ExecuteTemplate(w, name, ctx)
		if err != nil {
			panic(err)
		}
	}
}

func Serve() error {
	r := chi.NewRouter()
	r.Get("/", t("index.html").Serve(func(ctx *Ctx) error {
		return errors.New("fak")
	}))

	r.Get("/entities", t("entities.html").Serve(func(ctx *Ctx) error {
		{
			in := entities.ListRequest{
				PaginationToken: u.Ptr(ctx.URL.Query().Get("paginationToken")),
			}
			if urltype := ctx.URL.Query().Get("type"); urltype != "" {
				in.Type = &urltype
			}

			res, err := entities.List(ctx, in)
			if err != nil {
				return err
			}
			ctx.Data["Entities"] = res.Data
			ctx.Data["PaginationToken"] = res.PaginationToken

			ctx.ResponseWriter.Header().Set("HX-Push-Url", ctx.URL.RequestURI())
		}
		{
			res, err := entities.ListTypes(ctx)
			if err != nil {
				return err
			}
			ctx.Data["EntityTypes"] = res.Data
		}

		return nil
	}))

	r.HandleFunc("/entities:create", t("entities_create.html").Serve(func(ctx *Ctx) error {
		if ctx.Method == "POST" {
			ctx.ParseForm()
			id := ctx.Form.Get("id")
			typ := ctx.Form.Get("type")
			e, err := entities.Create(ctx, entities.CreateRequest{
				ID:   id,
				Type: typ,
			})
			ctx.Data["Entity"] = e
			fmt.Println(e)
			ctx.ResponseWriter.Header().Set("HX-Redirect", "/entities/"+e.Type+":"+e.ID)
			return err
		}
		// res, err := entities.List(ctx, entities.ListRequest{
		// })
		// if err != nil {
		// 	return err
		// }
		//
		// ctx.Data["Entities"] = res.Data
		// ctx.Data["PaginationToken"] = res.PaginationToken
		return nil
	}))

	r.Handle("/entities/{type}:{id}", t("entities_get.html").Serve(func(ctx *Ctx) error {
		typ, id := chi.URLParamFromCtx(ctx, "type"), chi.URLParamFromCtx(ctx, "id")
		e, err := entities.Get(ctx, id, typ)
		if err != nil {
			return err
		}

		ctx.Data["Entity"] = e
		return nil
	}))

	return http.ListenAndServe(":8000", r)
}
