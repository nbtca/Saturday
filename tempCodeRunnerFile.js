function testError() {
  try {
    throw new Error("Error");
  } catch (error) {
    return error;
  }
}
function wrapper() {
  try {
    testError();
  } catch (error) {
    return error;
  }
}
try {
  wrapper();
} catch (error) {
  console.log(error);
}
