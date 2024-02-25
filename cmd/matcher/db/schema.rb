# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[7.1].define(version: 2024_02_24_055906) do
  create_table "clients", force: :cascade do |t|
    t.string "client_id"
    t.string "session_id"
    t.string "client_key"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["client_id"], name: "index_clients_on_client_id", unique: true
  end

  create_table "sessions", force: :cascade do |t|
    t.string "session_id"
    t.string "router_addr"
    t.string "owner_id"
    t.string "guest_id"
    t.string "owner_client_id"
    t.string "guest_client_id"
    t.string "name", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["owner_id", "guest_id"], name: "index_sessions_on_owner_id_and_guest_id"
    t.index ["session_id"], name: "index_sessions_on_session_id", unique: true
  end

  create_table "users", force: :cascade do |t|
    t.string "user_id"
    t.string "name"
    t.string "login_id"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["login_id"], name: "index_users_on_login_id", unique: true
    t.index ["user_id"], name: "index_users_on_user_id", unique: true
  end

end
