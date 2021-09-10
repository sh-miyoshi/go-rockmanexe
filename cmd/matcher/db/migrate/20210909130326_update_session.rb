class UpdateSession < ActiveRecord::Migration[6.1]
  def change
    rename_column :sessions, :client_id, :owner_client_id
    rename_column :sessions, :client_key, :owner_client_key
    add_column :sessions, :guest_client_id, :string
    add_column :sessions, :guest_client_key, :string
    add_column :sessions, :route_id, :string
  end
end
