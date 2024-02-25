class CreateClients < ActiveRecord::Migration[7.1]
  def change
    create_table :clients do |t|
      t.string :client_id
      t.string :session_id
      t.string :client_key

      t.timestamps
    end

    add_index :clients, :client_id, unique: true
  end
end
