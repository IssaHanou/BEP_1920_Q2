## HAO video presentation

1. check IP: WiFi -> properties -> IPv4 
2. put IP into `back-end/resources/production/issa_config.json`
3. put IP into `broker/me.conf`
4. put IP into `cc_library/issa/slides/src/assets/slides_config.json`
5. put IP into `cc_library/issa/navigate/src/assets/navigate_config.json`
6. start mosquitto server using `me.conf`
7. start back-end with `issa_config.json`
8. start slides: 
    * `cd cc_library/issa/slides`
    * `ng serve`
9. start navigate:
    * `cd cc_library/issa/navigate`
    * `ng serve --host 192.168.0.107`
