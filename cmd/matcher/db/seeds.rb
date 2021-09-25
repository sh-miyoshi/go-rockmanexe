require "securerandom"

# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).

if Rails.env == "development"
  # Create test session
  session = Session.new(
      {
        session_name: SecureRandom.uuid,
        router_addr: "localhost:80",
        owner_id: SecureRandom.uuid,
        owner_client_id: "tester1",
        owner_client_key: "testtest",
        guest_id: SecureRandom.uuid,
        guest_client_id: "tester2",
        guest_client_key: "testtest",
        session_id: SecureRandom.uuid
      }
    )

    Client.create(
      {
        client_id: session.owner_client_id,
        client_key: session.owner_client_key,
        session_id: session.session_id
      }
    )
    Client.create(
      {
        client_id: session.guest_client_id,
        client_key: session.guest_client_key,
        session_id: session.session_id
      }
    )

    session.save!
end
