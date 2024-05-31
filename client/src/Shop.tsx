import React, { useState, useEffect } from 'react';
import axios from 'axios';

interface Item {
  id: string;
  name: string;
  description: string;
  price: number;
}

const Shop: React.FC<{ user: any; onLogout: () => void }> = ({ user, onLogout }) => {
  const [items, setItems] = useState<Item[]>([]);
  const [cart, setCart] = useState<Item[]>([]);

  useEffect(() => {
    const fetchItems = async () => {
      try {
        const response = await axios.get(`${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/items`);
        setItems(response.data);
      } catch (error) {
        console.error('ERROR: could not fetch items:', error);
      }
    };
    fetchItems();
  }, []);

  // append item to list of items stored in the cart
  const addToCart = (item: Item) => {
    setCart([...cart, item]);
  };

  // remove item from list (cart)
  const removeCartItem = async (id: string) => {
    try {
      await axios.delete(`/items/${id}`);
      setItems(items.filter(item => item.id !== id));
    } catch (error) {
      console.error('ERROR: could not remove item from cart:', error);
    }
  };

  // finalizing order, i.e., "buying" items within cart
  const handleSubmit = async () => {
    try {
      const orderDetails = cart.map(item => `${item.name}: ${item.price}`).join('\n');
      const message = `Order details:\n${orderDetails}`;
      await axios.post('/send', {
        name: user.name,
        email: user.email,
        message,
        cart
      });
      alert('Your order has successfully been sent for processing!');
    } catch (error) {
      console.error('ERROR: could not send order:', error);
      alert('ERROR: There was an issue sending your order. Please try again');
    }
  };

  const handleLogout = () => {
    onLogout(); 
  };

  return (
    <div className="Shop">
      {user && <button onClick={handleLogout}>Logout</button>}
      <h2>Items</h2>
    </div>
  );
};

export default Shop;
