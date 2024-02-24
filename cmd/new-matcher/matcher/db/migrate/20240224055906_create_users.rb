class CreateUsers < ActiveRecord::Migration[7.1]
  def change
    create_table :users, id: :uuid do |t|
      t.string :name
      t.string :login_id

      t.timestamps
    end

    add_index :users, :login_id, unique: true
  end
end
