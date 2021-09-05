class UpdateSessions < ActiveRecord::Migration[6.1]
  def change
    rename_column :sessions, :owner, :owner_id
    add_column :sessions, :guest_id, :string
    add_column :sessions, :session_name, :string

    add_index :sessions, [:owner_id, :guest_id]
    change_column :sessions, :session_name, :string, null: false
  end
end
