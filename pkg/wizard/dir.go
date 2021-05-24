package wizard

func getObjectMap() []*Object {
	return []*Object{
		{
			Name: NamePlaceholder,
			Type: TypeDir,
			SubObjects: []*Object{
				{
					Name: "api",
					Type: TypeDir,
					SubObjects: []*Object{
						{
							Name: "openapi",
							Type: TypeDir,
						},
						{
							Name: "proto",
							Type: TypeDir,
						},
					},
				},
				{
					Name: "application",
					Type: TypeDir,
				},
				{
					Name: "cmd",
					Type: TypeDir,
					SubObjects: []*Object{
						{
							Name: NamePlaceholder,
							Type: TypeDir,
							SubObjects: []*Object{
								{
									Name:     "main.go",
									Type:     TypeFile,
									Template: getTemplate("main.go"),
								},
								{
									Name:     "main_test.go",
									Type:     TypeFile,
									Template: getTemplate("main_test.go"),
								},
							},
						},
					},
				},
				{
					Name: "docker",
					Type: TypeDir,
					SubObjects: []*Object{
						{
							Name:     "Dockerfile",
							Type:     TypeFile,
							Template: getTemplate("Dockerfile"),
						},
					},
				},
				{
					Name: "domain",
					Type: TypeDir,
				},
				{
					Name: "infrastructure",
					Type: TypeDir,
					SubObjects: []*Object{
						{
							Name: "transport",
							Type: TypeDir,
							SubObjects: []*Object{
								{
									Name: "http",
									Type: TypeDir,
								},
								{
									Name: "grpc",
									Type: TypeDir,
								},
								{
									Name: "amqp",
									Type: TypeDir,
								},
							},
						},
						{
							Name: "repositories",
							Type: TypeDir,
						},
					},
				},
				{
					Name: "logs",
					Type: TypeDir,
				},
				{
					Name:     "docker-compose.yml",
					Type:     TypeFile,
					Template: getTemplate("docker-compose.yml"),
				},
				{
					Name:     "Dockerfile",
					Type:     TypeFile,
					Template: getTemplate("Dockerfile.dev"),
				},
				{
					Name:     "go.mod",
					Type:     TypeFile,
					Template: getTemplate("go.mod"),
				},
				{
					Name:     ".env",
					Type:     TypeFile,
					Template: getTemplate(".env"),
				},
				{
					Name:     ".env.example",
					Type:     TypeFile,
					Template: getTemplate(".env"),
				},
				{
					Name:     ".gitignore",
					Type:     TypeFile,
					Template: getTemplate(".gitignore"),
				},
			},
		},
	}
}
