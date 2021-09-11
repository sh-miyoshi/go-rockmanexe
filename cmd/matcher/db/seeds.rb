require "securerandom"

# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).

user_id = SecureRandom.uuid
User.create(name: "ロックマン", login_id: "login_id_1", user_id: user_id)
user2_id = SecureRandom.uuid
User.create(name: "ブルース", login_id: "login_id_2", user_id: user2_id)

Session.create(
  session_name: "ライバル対戦",
  owner_id: user_id,
  router_addr: Settings.router.data_addr,
  owner_client_id: "tester1",
  owner_client_key: "testtest",
  expires_at: Time.local(2021, 12, 31),
  guest_id: user2_id,
  guest_client_id: "tester2",
  guest_client_key: "testtest",
  route_id: "nagnlabmklabjng"
)

