# NB The lurl app runs in http because it expects any loadbalancing and SSL to
# be handled upstream
# As such, there will only be http assertions in the test suite

# A known key
GET  http://localhost:8080/gh
HTTP 302
[Asserts]
header "Location" matches "https://www.github.com/"

# Another known key
GET  http://localhost:8080/e
HTTP 302
[Asserts]
header "Location" matches "https://www.example.com/"

# A known missing key
GET  http://localhost:8080/404
HTTP 404

