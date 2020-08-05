# file: features/image_previewer.feature

Feature: Image previewer in work
	Scenario: The image found in cache
		When I send "GET" request to "http://image_previewer:6080/500/200/nginx/image.jpg" and send "GET" second request to "http://image_previewer:6080/500/200/nginx/image.jpg"
		Then The response code should be 200 and resonse timeout less "10ms"

	Scenario: Remote server not found
		When I send "GET" request to "http://image_previewer:6080/500/200/foo.bar/image.jpg"
		Then The response code should be 500

	Scenario: Remote server is found, but image not found
		When I send "GET" request to "http://image_previewer:6080/500/200/nginx/not_found_image.jpg"
		Then The response code should be 502

	Scenario: Remote server is found, but image file is not image file
		When I send "GET" request to "http://image_previewer:6080/500/200/nginx/libdbm64.so"
		Then The response code should be 500

	Scenario: Remote server return error
		When I send "GET" request to "http://image_previewer:6080/500/200/nginx/"
		Then The response code should be 502

	Scenario: Remote server return image
		When I send "GET" request to "http://image_previewer:6080/100/200/nginx/image.jpg"
		Then The response code should be 200

	Scenario: Image size too small
		When I send "GET" request to "http://image_previewer:6080/2285/1609/nginx/image.jpg"
		Then The response code should be 200

	Scenario: Input image size too large
		When I send "GET" request to "http://image_previewer:6080/200/100/nginx/large_image.jpg"
		Then The response code should be 500

	Scenario: Output image size too large
		When I send "GET" request to "http://image_previewer:6080/3428/2414/nginx/image.jpg"
		Then The response code should be 400

	Scenario: Input image content too large
		When I send "GET" request to "http://image_previewer:6080/200/100/nginx/super_large_image.jpg"
		Then The response code should be 500
