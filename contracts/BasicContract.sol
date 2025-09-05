pragma solidity ^0.5.0;

contract BasicContract {
    uint256 public value;
    
    constructor() public {
        value = 1000;
    }
    
    function get() public view returns (uint256) {
        return value;
    }
    
    function set(uint256 _value) public {
        value = _value;
    }
}
