# file: features/image_previewer.feature

# http://localhost:6080/
# http://image_previewer:6080/

Feature: Image previewer in work
	Scenario: The image found in cache
		When I send "GET" request to "http://image_previewer:6080/500/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg" and send "GET" second request to "http://image_previewer:6080/500/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg"
		Then The response code should be 200 and resonse timeout less "10ms"

	Scenario: Remote server not found
		When I send "GET" request to "http://image_previewer:6080/500/200/foo.bar/some_image.jpg"
		Then The response code should be 500

	Scenario: Remote server is found, but image not found
		When I send "GET" request to "http://image_previewer:6080/500/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk123.jpga1_1902_16_barred-owl_sandra_rothenberg_kk123.jpg"
		Then The response code should be 502

	Scenario: Remote server is found, but image file is not image file
		When I send "GET" request to "http://image_previewer:6080/500/200/www.audubon.org/"
		Then The response code should be 500

	Scenario: Remote server return error
		When I send "GET" request to "http://image_previewer:6080/500/200/www.audubon.org/sites/default/files/"
		Then The response code should be 502

	Scenario: Remote server return image
		When I send "GET" request to "http://image_previewer:6080/100/200/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg"
		Then The response code should be 200

	Scenario: Image size to small
		When I send "GET" request to "http://image_previewer:6080/3428/2414/www.audubon.org/sites/default/files/a1_1902_16_barred-owl_sandra_rothenberg_kk.jpg"
		Then The response code should be 200
