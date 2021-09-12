class DeleteExpiresAt < ActiveRecord::Migration[6.1]
  def change
    remove_column :sessions, :expires_at, :date
  end
end
