module ApiHelper
  extend ActiveSupport::Concern

  private

  # API status responder
  def response_bad_request
    render status: :bad_request, json: { status: 400, message: "Bad Request" }
  end

  def response_unauthorized
    render status: :unauthorized, json: { status: 401, message: "Unauthorized" }
  end

  def response_not_found(class_name: "page")
    render status: :not_found, json: { status: 404, message: "#{class_name.capitalize} Not Found" }
  end

  def response_internal_server_error
    render status: :internal_server_error, json: { status: 500, message: "Internal Server Error" }
  end
end
