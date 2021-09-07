class AddUniqueToId < ActiveRecord::Migration[6.1]
  def change
    add_index :users, :user_id, unique: true
    add_index :users, :login_id, unique: true
  end
end
