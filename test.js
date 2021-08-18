function testError() {
  throw new Error("Error");
}
function wrapper() {
  testError();
}
try {
  wrapper();
} catch (error) {
  console.log(error);
}
