require "securerandom"

# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).
User.create(name: "tester1", login_id: "login_id_1", user_id: SecureRandom.uuid)
