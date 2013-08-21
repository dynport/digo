# digo

[![Build Status](https://travis-ci.org/dynport/digo.png?branch=master)](https://travis-ci.org/dynport/digo)

DigitalOcean API cli tool and library for golang

## Requirements

Set the following env variables:
    
    export DIGITAL_OCEAN_CLIENT_ID=<secret>
    export DIGITAL_OCEAN_API_KEY=<secret>

These env settings can be set optionally:

    export DIGITAL_OCEAN_DEFAULT_REGION_ID=2
    export DIGITAL_OCEAN_DEFAULT_SIZE_ID=66
    export DIGITAL_OCEAN_DEFAULT_IMAGE_ID=350076
    export DIGITAL_OCEAN_DEFULT_SSH_KEY=<secret_id>

## Installation
    
    go get -u github.com/dynport/digo

## Usage
    $ digo --help
    USAGE
    droplet	create 	<name>      	Create new droplet                            
                                  -i DEFAULT: "350076" Image id for new droplet 
                                  -r DEFAULT: "2"      Region id for new droplet
                                  -s DEFAULT: "66"     Size id for new droplet  
                                  -k DEFAULT: "<secret>"  Ssh key to be used       
    droplet	destroy	<droplet_id>	Destroy droplet                               
    droplet	list   	            	List active droplets                          
    droplet	rebuild	<droplet_id>	Rebuild droplet                               
                                  -i DEFAULT: "0" Rebuild droplet               
    key    	list   	            	List available ssh keys                       
    region 	list   	            	List available droplet regions                
    size   	list   	            	List available droplet sizes                  
