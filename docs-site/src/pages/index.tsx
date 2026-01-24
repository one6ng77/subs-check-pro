import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import styles from './index.module.css';
import HomepageFeatures from '../components/HomepageFeatures';

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout title={siteConfig.title} description="测活、测速、媒体解锁，代理检测工具，自动生成 mihomo 和 singbox 订阅">
      {/* 使用 homeWrapper 包装，实现垂直居中布局 */}
      <div className={styles.homeWrapper}>
        <main className={styles.mainContent}>
          <div className={styles.heroContainer}>
            <h1 className={styles.title}>{siteConfig.title}</h1>
            <p className={styles.subtitle}>
              测活、测速、媒体解锁，代理检测工具，自动生成 mihomo 和 singbox 订阅
            </p>

            <div className={styles.actions}>
              <Link className={styles.primaryBtn} to="/docs/Home">
                使用说明
              </Link>
            </div>
          </div>
          
          <HomepageFeatures/>
        </main>
      </div>
    </Layout>
  );
}