class CreateSessions < ActiveRecord::Migration[6.1]
  def change
    create_table :sessions do |t|
      t.string :owner
      t.string :router_addr
      t.string :client_id
      t.string :client_key
      t.date :expires_at

      t.timestamps
    end
  end
end
