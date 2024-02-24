class CreateClients < ActiveRecord::Migration[7.1]
  def change
    create_table :clients, id: :uuid do |t|
      t.string :session_id
      t.string :key

      t.timestamps
    end
  end
end
