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

ActiveRecord::Schema.define(version: 2021_09_05_124113) do

  create_table "histories", force: :cascade do |t|
    t.text "users"
    t.date "finished_at"
    t.string "winner"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "sessions", force: :cascade do |t|
    t.string "owner_id"
    t.string "router_addr"
    t.string "client_id"
    t.string "client_key"
    t.date "expires_at"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.string "guest_id"
    t.string "session_name", null: false
    t.index ["owner_id", "guest_id"], name: "index_sessions_on_owner_id_and_guest_id"
  end

  create_table "users", force: :cascade do |t|
    t.string "name"
    t.string "login_id"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.string "user_id", null: false
  end

end
