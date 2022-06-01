package hypergo

import "fmt"

func Example() {
	head := Head().Title().Text("Title")

	span := Span().Text("Hello World!")

	ul := Ul()
	for i := 0; i < 3; i++ {
		ul.Li().Textv("Count: ", i)
	}

	btn := Button("type=button", "id=foo", "disabled").Text("Click me!")
	form := Form().Append(
		Label().Append(
			Text("input number"),
			Input("type=number", "value=1"),
		),
		btn,
	)

	body := Body().Append(
		span,
		ul,
		form,
	)

	doc := HTML5(Html().Append(head, body))

	fmt.Println(doc.RenderIndent("\t"))

	// Output:
	// <!DOCTYPE html>
	// <html>
	// 	<head>
	// 		<title>Title</title>
	// 	</head>
	// 	<body>
	// 		<span>Hello World!</span>
	// 		<ul>
	// 			<li>Count: 0</li>
	// 			<li>Count: 1</li>
	// 			<li>Count: 2</li>
	// 		</ul>
	// 		<form>
	// 			<label>
	// 				input number
	// 				<input type="number" value="1">
	// 			</label>
	// 			<button type="button" id="foo" disabled>Click me!</button>
	// 		</form>
	// 	</body>
	// </html>
}
