import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Popup from './Popup';

// Mock chrome API
const chromeMock = {
  storage: {
    local: {
      get: vi.fn(),
      set: vi.fn(),
    },
  },
  tabs: {
    query: vi.fn(),
    reload: vi.fn(),
  },
};

global.chrome = chromeMock as unknown as typeof chrome;

describe('Popup', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders correctly', () => {
    // Mock storage get to return default
    chromeMock.storage.local.get.mockImplementation((keys, callback) => {
      callback({ enabled: true });
    });

    render(<Popup />);
    
    expect(screen.getByText('YouTube Replacer')).toBeInTheDocument();
    expect(screen.getByText('Active')).toBeInTheDocument();
  });

  it('toggles switch', () => {
    chromeMock.storage.local.get.mockImplementation((keys, callback) => {
      callback({ enabled: true });
    });

    render(<Popup />);
    
    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);
    
    expect(chromeMock.storage.local.set).toHaveBeenCalledWith({ enabled: false });
  });
});
