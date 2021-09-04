require "test_helper"

class LoginControllerTest < ActionDispatch::IntegrationTest
  test "should get create" do
    get login_create_url
    assert_response :success
  end

  test "should get destroy" do
    get login_destroy_url
    assert_response :success
  end

  test "should get callback" do
    get login_callback_url
    assert_response :success
  end
end
