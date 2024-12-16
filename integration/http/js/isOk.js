client.test("Test status code", function() {
    console.log("Response status: " + response.status);
    client.assert(response.status === 200, "Response status is not 200");
});