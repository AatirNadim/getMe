package net.getMeStore.client.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "getme.store")
public class Env {
    private String socketPath;

    public void setSocketPath(String socketPath) {
        this.socketPath = socketPath;
    }


    public String getSocketPath() {
        return this.socketPath;
    }
}
