package net.getMeStore.client.config;

import net.getMeStore.client.service.GetMeService;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
@EnableConfigurationProperties(Env.class)
public class GetMeAutoConfiguration {

    @Bean
    public UdsHandler udsHandler() {
        return new UdsHandler();
    }

    @Bean
    public GetMeService getMeService(UdsHandler udsHandler, Env env) {
        return new GetMeService(udsHandler, env);
    }
}
