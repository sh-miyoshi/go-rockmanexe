class AddClientKey < ActiveRecord::Migration[6.1]
  def change
    add_column :clients, :client_key, :string
  end
end
