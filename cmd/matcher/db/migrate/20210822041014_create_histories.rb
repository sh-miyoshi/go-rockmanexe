class CreateHistories < ActiveRecord::Migration[6.1]
  def change
    create_table :histories do |t|
      t.text :users, array: true
      t.date :finished_at
      t.string :winner

      t.timestamps
    end
  end
end
