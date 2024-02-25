class CreateUsers < ActiveRecord::Migration[7.1]
  def change
    create_table :users do |t|
      t.string :user_id
      t.string :name
      t.string :login_id

      t.timestamps
    end

    add_index :users, :user_id, unique: true
    add_index :users, :login_id, unique: true
  end
end
