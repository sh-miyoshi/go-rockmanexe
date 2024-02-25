class Session < ApplicationRecord
  validates :name, length: { in: 3..20 }
  validates :owner_id, presence: true
  validates :guest_id, presence: true
end
