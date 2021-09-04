class LoginController < ApplicationController
  def create
    uri = URI.parse("#{Settings.login[:server_addr]}")

    state = SecureRandom.hex(12)
    verifier = SecureRandom.hex(128)
    challenge = get_code_challenge(verifier)
    redirect_uri = "#{Settings.login[:bbs_addr]}/login/callback"

    queries = {
      'scope' => 'openid email',
      'response_type' => 'code',
      'client_id' => Settings.login[:client_id],
      'client_secret' => Settings.login[:client_secret],
      'redirect_uri' => redirect_uri,
      'code_challenge' => challenge,
      'code_challenge_method' => 'S256',
      'state' => state
    }
    uri.query = URI.encode_www_form(queries)

    session[:state] = state
    session[:verifier] = verifier

    logger.debug("login redirect to #{uri}")
    redirect_to uri.to_s
  end

  def destroy
  end

  def callback
  end

  private

  def get_code_challenge(verifier)
    # currentryl supported only S256
    digest = OpenSSL::Digest.new('sha256')
    Base64.urlsafe_encode64(digest.update(verifier).digest).delete('=')
  end
end
