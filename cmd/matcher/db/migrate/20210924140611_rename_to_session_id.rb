class RenameToSessionId < ActiveRecord::Migration[6.1]
  def change
    rename_column :sessions, :route_id, :session_id
  end
end
