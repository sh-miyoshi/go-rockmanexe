class ApplicationController < ActionController::Base
  private

  # API status responder

  def response_bad_request
    render status: 400, json: { status: 400, message: 'Bad Request' }
  end

  def response_unauthorized
    render status: 401, json: { status: 401, message: 'Unauthorized' }
  end

  def response_not_found(class_name: 'page')
    render status: 404, json: { status: 404, message: "#{class_name.capitalize} Not Found" }
  end

  def response_internal_server_error
    render status: 500, json: { status: 500, message: 'Internal Server Error' }
  end
end
