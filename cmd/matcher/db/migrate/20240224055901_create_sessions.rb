class CreateSessions < ActiveRecord::Migration[7.1]
  def change
    create_table :sessions do |t|
      t.string :session_id
      t.string :router_addr
      t.string :owner_id
      t.string :guest_id
      t.string :owner_client_id
      t.string :guest_client_id
      t.string :name, null: false

      t.timestamps
    end

    add_index :sessions, :session_id, unique: true
    add_index :sessions, %i[owner_id guest_id]
  end
end
